import { useState } from "react"
import Link from "next/link"
import { useSession, getSession } from "next-auth/client"
import Router from "next/router"

////import axios from 'axios'
//import ProtectedRoute, { useAuth } from "../context/auth";
//import CreateRoom from "../components/CreateRoomModal";

const fetchData = async () => {
  //console.log(ctx);
  const res = await fetch("https://chat.tony.place/get")
  const data = await res.json()

  if (!data) {
    return {
      notFound: true,
    }
  }

  return data
}

export async function getServerSideProps(ctx) {
  return { props: { data: await fetchData(), session: await getSession(ctx) } }
}

function Home({ data }) {
  const [session, loading] = useSession()
  const [rooms, setRooms] = useState(data?.rooms)
  const [showModal, setShowModal] = useState(false)

  if (!session && typeof window !== "undefined") {
    window.location = "/login"
    return null
  }

  const createRoom = async () => {
    try {
      let temp = rooms || []
      //Switch /create to a POST and take room data from the frontend
      const res = await fetch("https://chat.tony.place/create")
      const data = await res.json()
      console.log(data)
      temp.push(data.roomID)
      setRooms([...temp])
    } catch (error) {
      console.log("Error Creating Room: ", error)
    }
    //    setShowModal(true)
  }

  return (
    <>
      {/*<CreateRoom show={showModal} setShow={setShowModal} />*/}
      <div className="w-1/2 mx-auto mt-10 text-center">
        <h1 className="text-4xl"> Voice Chat </h1>

        <div className="flex flex-row justify-between mt-10">
          <div>
            <button
              className="flex-grow-0 px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white"
              onClick={createRoom}
            >
              {" "}
              Create Room{" "}
            </button>
          </div>
          <div className="w-3/4 flex-grow-1">
            <span className="block pb-4"> Open Rooms </span>
            <ul>
              {" "}
              {rooms?.length > 0 ? (
                rooms?.map((e, k) => {
                  return (
                    <li key={k}>
                      <Link href={`/Room/${encodeURIComponent(e?.roomUuid)}`}>
                        <a className="block px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white">
                          {" "}
                          {e?.roomUuid}{" "}
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
    </>
  )
}

export default Home
