import { useEffect, useState } from "react"
import Link from "next/link"
import nookies from "nookies"
import { useRouter } from "next/router"
import { useAuth } from "../context/auth"

const fetchData = async (ctx) => {
  try {
    const cookies = nookies.get(ctx)
    let headers = new Headers()
    headers.append("Authorization", `Bearer ${cookies.token}`)
    const res = await fetch("http://localhost:8000/get", { headers })
    const data = await res.json()

    if (!data) {
      return {
        notFound: true,
      }
    }

    return data
  } catch (error) {
    console.log("err", error)
    return {
      notFound: true,
    }
  }
}

export async function getServerSideProps(ctx) {
  return {
    props: {
      data: await fetchData(ctx),
    },
  }
}

function Home({ data }) {
  const router = useRouter()
  const { user, loading, signOut } = useAuth()
  const [rooms, setRooms] = useState(data?.rooms)

  useEffect(() => {
    if (!loading && !user) {
      console.log("login component user: ", user)
      router.push("/login")
    }
  }, [loading])

  const handleCreateRoom = async () => {
    //try {
    //  let temp = rooms || []
    //  //Switch /create to a POST and take room data from the frontend
    //  const res = await fetch("https://chat.tony.place/create")
    //  const data = await res.json()
    //  console.log(data)
    //  temp.push(data.roomID)
    //  setRooms([...temp])
    //} catch (error) {
    //  console.log("Error Creating Room: ", error)
    //}
    //    setShowModal(true)
  }

  return !loading && user ? (
    <>
      {/*<CreateRoom show={showModal} setShow={setShowModal} />*/}
      <div className="w-1/2 mx-auto mt-10 text-center">
        <h1 className="text-4xl"> Voice Chat </h1>
        <button
          onClick={async () => {
            const isSignedOut = await signOut()
            if (isSignedOut) {
              router.replace("/login")
            }
          }}
        >
          Sign out
        </button>
        <div className="flex flex-row justify-between mt-10">
          <div>
            <button
              className="flex-grow-0 px-2 py-2 border-2 rounded hover:bg-blue-700 hover:text-white"
              onClick={handleCreateRoom}
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
  ) : (
    <> </>
  )
}

export default Home
