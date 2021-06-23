import 'tailwindcss/tailwind.css'

function MyApp({ Component, pageProps }) {
  return (
    <>
      <style jsx global>
        {`
          body {
            background: linear-gradient(
                70.12deg,
                rgba(0, 147, 233, 0.5) 0%,
                rgba(128, 208, 199, 0.5) 100%
              ),
              #ffffff;
          }
          .box-shadow {
            box-shadow: 0px 0px 5px 0.5px rgba(0, 0, 0, 0.16);
          }
        `}
      </style>
      <Component {...pageProps} />
    </>
  )
}

export default MyApp
