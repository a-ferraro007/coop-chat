import { useEffect } from 'react'

function CreateRoom(props) {
  useEffect(() => {
    console.log(props)
  }, [props])

  const createRoom = async () => {
    try {
      let temp = rooms || []
      //Switch /create to a POST and take room data from the frontend
      const res = await axios.get('http://localhost:8000/create')
      console.log(res.data.roomUuid)
      temp.push(res.data.roomUuid)
      setRooms([...temp])
    } catch (error) {
      console.log('Error Creating Room: ', error)
    }
  }
  return (
    <>
      {props.show ? (
        <div className="absolute top-0 bottom-0 left-0 right-0 w-64 h-52 mx-auto bg-bg-white">
          <span>Create A Room</span>
          <button
            onClick={() => {
              props.setShow(false)
            }}
          >
            {' '}
            Close{' '}
          </button>{' '}
        </div>
      ) : (
        <></>
      )}
    </>
  )
}

export default CreateRoom
