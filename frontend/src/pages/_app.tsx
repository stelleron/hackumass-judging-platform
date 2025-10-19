// pages/_app.tsx
import Head from "next/head";
import "../styles/globals.css";
import type { AppProps } from "next/app";

export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <title>HackUMass XIII Judging</title>
      </Head>
      <div>
        <Component {...pageProps} />
      </div>
    </>
  );
}
