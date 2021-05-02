import { useState, useEffect } from 'react'
import Link from 'next/link'

export async function getStaticProps() {
  const asyncGet = async () => {
    let res = await fetch('http://localhost:8000/get', {
      method: 'GET'
    })
    //console.log(res)
    return res.json()
  }

  try {
    const data = await asyncGet()
    console.log(data)
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

  const createCoop = () => {
    const asyncGet = async () => {
      let res = await fetch('http://localhost:8000/create', {
        method: 'GET'
      })
      return res.json()
    }

    asyncGet().then((data) => {
      let tmp = rooms
      tmp.push(data.roomID)
      console.log(tmp)
      setRooms([...tmp])
    })

    //console.log(res.json().then((data) => console.log(data)))
  }

  useEffect(() => {
    console.log(rooms)
  }, [])

  //useEffect(async () => {
  //  const getData = async () => {
  //    const asyncGet = async () => {
  //      let res = await fetch('http://localhost:8000/get', {
  //        method: 'GET'
  //      })
  //      console.log(res)
  //      return res.json()
  //    }

  //    asyncGet().then((data) => {
  //      //tmp.push(data.roomID)
  //      //console.log(tmp)
  //      //let tmp = rooms
  //      //tmp.push(data.rooms)
  //      console.log(data)
  //      setRooms([...data.rooms])
  //    })
  //  }
  //  await getData()
  //}, [])

  return (
    <div className="w-1/2 mx-auto text-center">
      <h1 className="text-4xl"> Coop Chat </h1>

      <div className="flex flex-row justify-between mt-10">
        <div>
          <button
            className="border-2 rounded px-2 py-2 hover:bg-blue-700 hover:text-white flex-grow-0"
            onClick={createCoop}
          >
            {' '}
            Create Room{' '}
          </button>
        </div>
        <div className="w-3/4 flex-grow-1">
          <span className="pb-4 block"> Open Rooms </span>
          <ul>
            {' '}
            {rooms?.length > 0 ? (
              rooms?.map((e, k) => {
                console.log(e)
                return (
                  <li key={k}>
                    <Link href={`/Room/${encodeURIComponent(e)}`}>
                      <a className=" block border-2 rounded px-2 py-2 hover:bg-blue-700 hover:text-white">
                        {' '}
                        {e}{' '}
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
