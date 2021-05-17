import { useRouter } from 'next/router'
import { useEffect, useState, useRef } from 'react'
import PeerVideo from '../../components/PeerVideo'

const Room = () => {
  const router = useRouter()
  const ws = useRef()
  const peers = useRef(new Object())
  const roomCount = useRef() //using a ref to avoid rerenders when updating roomSize in the WS
  const audioStream = useRef()
  const userAudio = useRef()
  const uuidRef = useRef('')
  const [peerStreams, setPeerStreams] = useState([])
  const [stateUUID, setUUID] = useState('')
  const [roomID, setRoomID] = useState(router.query.roomID)

  useEffect(() => {
    if (router.isReady) {
      setRoomID(router.query.roomID)
    }
  }, [router])

  const getMedia = async () => {
    const constraints = {
      audio: { echoCancellation: true, noiseSuppression: true },
      video: true
    }
    if (navigator.mediaDevices) {
      try {
        return await navigator.mediaDevices.getUserMedia(constraints)
      } catch (err) {
        console.log('Error accessing media devices: ', err)
      }
    }
  }

  useEffect(() => {
    //Having a ref and a usestate for uuid doesn't make sense. Fix this.
    uuidRef.current = createUUID()
    setUUID(uuidRef.current)
    console.log('UUID', uuidRef.current)
    getMedia().then((stream) => {
      userAudio.current.srcObject = stream
      console.log(userAudio.current)
      audioStream.current = stream

      ws.current = new WebSocket(`ws://localhost:8000/join?roomID=${roomID}`)
      ws.current.onopen = () => {
        console.log('Open WS Connection ...')
        ws.current.send(
          JSON.stringify({
            join: true,
            uuid: uuidRef.current,
            dest: 'pool'
          })
        )
      }

      ws.current.onmessage = (message) => {
        const data = JSON.parse(message.data)

        //If a user joins use callPeer to create new peerConnection and create offer
        //Update roomcount ref with new room size
        if (data.join && data.dest === 'pool') {
          roomCount.current = data.size
          createPeerConn(data.uuid)

          ws.current.send(
            JSON.stringify({
              join: true,
              uuid: uuidRef.current,
              dest: data.uuid
            })
          )
        } else if (data.join && data.dest === uuidRef.current) {
          createPeerConn(data.uuid, true)
        } else if (data.offer && data.dest === uuidRef.current) {
          console.log('offer: ', data.offer)
          handleOffer(data.uuid, data.offer)
        } else if (data.candidate) {
          peers.current[data.uuid].peer
            .addIceCandidate(new RTCIceCandidate(data.candidate))
            .catch((error) => console.log(error))
        }
      }
    })

    return () => {
      //Close WS when component unmounts
      console.log('...disconnecting ')
      ws.current.close()
    }
  }, [])

  const createPeerConn = async (uuid, call = false) => {
    const config = { iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] }
    const peer = new RTCPeerConnection(config)
    peers.current[uuid] = { peer }

    peers.current[uuid].peer.onicecandidate = (event) =>
      handleOnIceCandidate(event, uuid)

    peers.current[uuid].peer.ontrack = (event) =>
      handleTrackEvent(event, peers.current[uuid].peer)
    audioStream.current.getTracks().forEach((track) => {
      peers.current[uuid].peer.addTrack(track, audioStream.current)
    })
    //peerConn.oniceconnectionstatechange disconnect and remove feed

    if (call) {
      try {
        const desc = await peers.current[uuid].peer.createOffer()
        createDescription(uuid, desc)
      } catch (error) {
        console.log('Error creating offer: ', error)
      }
    }
  }

  const createDescription = async (uuid, description) => {
    try {
      await peers.current[uuid].peer.setLocalDescription(description)
      ws.current.send(
        JSON.stringify({
          offer: peers.current[uuid].peer.localDescription,
          uuid: uuidRef.current,
          dest: uuid
        })
      )
    } catch (error) {
      console.log('Error Creating Description: ', error)
    }
  }

  const handleOnIceCandidate = (event, uuid) => {
    console.log('Found Ice Candidate')
    if (event.candidate) {
      ws.current.send(
        JSON.stringify({
          candidate: event.candidate,
          uuid: uuidRef.current,
          dest: uuid
        })
      )
    }
  }

  const handleTrackEvent = (event) => {
    if (event.track.kind === 'audio') return
    const temp = peerStreams
    temp.push(<PeerVideo key={event.streams[0].id} stream={event.streams[0]} />)
    setPeerStreams([...peerStreams])
    console.log('Track Event ', event.streams[0])
  }

  const handleOffer = async (uuid, offer) => {
    try {
      await peers.current[uuid].peer.setRemoteDescription(
        new RTCSessionDescription(offer)
      ) ///ERROR here on multiple participants

      if (offer.type === 'offer') {
        try {
          const description = await peers.current[uuid].peer.createAnswer()
          await createDescription(uuid, description)
        } catch (error) {
          console.log('Error Creating Answer: ', error)
        }
      }
    } catch (error) {
      console.log(error)
    }
  }

  // Taken from http://stackoverflow.com/a/105074/515584
  // Not a real UUID but works in this small scenario
  function createUUID() {
    function s4() {
      return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1)
    }

    return (
      s4() +
      s4() +
      '-' +
      s4() +
      '-' +
      s4() +
      '-' +
      s4() +
      '-' +
      s4() +
      s4() +
      s4()
    )
  }

  const handleLeaveRoom = () => {
    //Close getUserMedia tracks when a user leaves
    audioStream.current.getTracks().forEach((track) => {
      track.stop()
    })
    router.push('/')
  }
  return (
    <div className="w-1/2 mx-auto mt-10 text-center">
      {roomID ? (
        <>
          <h1 className="pb-10 text-4xl">
            {' '}
            Room: <span className="text-blue-700"> {roomID} </span>{' '}
          </h1>
          <h2>
            UUID: <span className="text-blue-700">{stateUUID}</span>
          </h2>
          <video
            muted
            id="audioDiv"
            autoPlay
            controls={true}
            ref={userAudio}
          ></video>
          <div>{peerStreams}</div>
          <button
            className="flex-grow-0 px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white"
            onClick={handleLeaveRoom}
          >
            {' '}
            Leave Room{' '}
          </button>{' '}
        </>
      ) : (
        <> </>
      )}
    </div>
  )
}

export default Room
