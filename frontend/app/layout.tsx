import type { Metadata } from "next";
import { JetBrains_Mono, Syne } from "next/font/google";
import "./globals.css";

const syne = Syne({
  variable: "--font-ui",
  subsets: ["latin"],
  weight: ["400", "500", "600"],
});

const jetBrainsMono = JetBrains_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
  weight: ["400", "500"],
});

export const metadata: Metadata = {
  title: "DOM Traverser - IF2211",
  description: "Interactive DOM traversal visualization for BFS, DFS, and LCA demo.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${syne.variable} ${jetBrainsMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}
