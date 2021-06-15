import { useState, useEffect } from 'react'
import Link from 'next/link'
import axios from 'axios'

export async function getStaticProps() {
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

export default function Home({ data }) {
  const [rooms, setRooms] = useState(data?.rooms)

  const createCoop = async () => {
    try {
      let temp = rooms || []
      //Switch /create to a POST and take room data from the frontend
      const res = await axios.get('http://localhost:8000/create')
      temp.push(res.data.roomID)
      setRooms([...temp])
    } catch (error) {
      console.log('Error Creating Room: ', error)
    }
  }

  return (
    <div className="w-1/2 mx-auto mt-10 text-center">
      <h1 className="text-4xl"> Coop Chat </h1>

      <div className="flex flex-row justify-between mt-10">
        <div>
          <button
            className="flex-grow-0 px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white"
            onClick={createCoop}
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
                console.log(e.roomUuid)
                return (
                  <li key={k}>
                    <Link href={`/Room/${encodeURIComponent(e.roomUuid)}`}>
                      <a className="block px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white">
                        {' '}
                        {e.roomUuid}{' '}
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
}
