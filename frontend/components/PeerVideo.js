import { useEffect, useRef } from 'react'

function PeerVideo(props) {
  const stream = useRef(props.stream)
  useEffect(() => {
    stream.current.srcObject = props.stream
  }, [])
  return (
    <>
      <video
        id="audioDiv"
        autoPlay
        controls={true}
        key={props.stream}
        ref={stream}
      ></video>
    </>
  )
}

export default PeerVideo
