// app/layout.tsx
import "../styles/globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "HackUMass XIII Judging Platform",
  description: "Login and manage judging for HackUMass XIII",
};

export default function RootLayout({ children, }: { children: React.ReactNode;}) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
