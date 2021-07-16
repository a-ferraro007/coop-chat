import { useState } from 'react'
import Link from 'next/link'
import axios from 'axios'
import { useAuth } from '../context/auth'
//import ProtectedRoute from '../context/protectedRoute'

import ProtectedRoute, { AuthProvider } from '../context/auth'

const fetchData = async () => {
  try {
    const res = await axios.get('http://localhost:8000/get')
    const data = res.data
    console.log(data.rooms[0])
    return {
      props: {
        data
      }
    }
  } catch (error) {
    return {
      props: {
        rooms: []
      }
    }
  }
}

export async function getStaticProps() {
  return await fetchData()
}

function Home({ data }) {
  const { user } = useAuth()
  const [rooms, setRooms] = useState(data?.rooms)

  const createRoom = async () => {
    try {
      let temp = rooms || []
      //Switch /create to a POST and take room data from the frontend
      const res = await axios.get('http://localhost:8000/create')
      console.log(res.data)
      temp.push(res.data.roomID)
      setRooms([...temp])
    } catch (error) {
      console.log('Error Creating Room: ', error)
    }
  }

  return (
    user && (
      <div className="w-1/2 mx-auto mt-10 text-center">
        <h1 className="text-4xl"> Voice Chat </h1>

        <div className="flex flex-row justify-between mt-10">
          <div>
            <button
              className="flex-grow-0 px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white"
              onClick={createRoom}
            >
              {' '}
              Create Room{' '}
            </button>
          </div>
          <div className="w-3/4 flex-grow-1">
            <span className="block pb-4"> Open Rooms </span>
            <ul>
              {' '}
              {rooms?.length > 0 ? (
                rooms?.map((e, k) => {
                  console.log(e?.roomUuid)
                  return (
                    <li key={k}>
                      <Link href={`/Room/${encodeURIComponent(e?.roomUuid)}`}>
                        <a className="block px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white">
                          {' '}
                          {e?.roomUuid}{' '}
                        </a>
                      </Link>
                    </li>
                  )
                })
              ) : (
                <div> No Rooms Found </div>
              )}
            </ul>
          </div>
        </div>
      </div>
    )
  )
}

export default ProtectedRoute(Home)
