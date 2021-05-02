import { route } from 'next/dist/next-server/server/router'
import { useRouter } from 'next/router'

const Room = () => {
  const router = useRouter()

  return <div>{router.query.roomID}</div>
}

export default Room
