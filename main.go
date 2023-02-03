package main

import(
    "log"
    "net"
    "fmt"
    "time"
    "io"
    "bytes"
    "encoding/binary"
    "crypto/rand"
)

type FileServer struct {}

func (fs *FileServer) start(){
     ln,err := net.Listen("tcp",":3000")
     if err != nil {
           log.Fatal(err)
     }
     for {
         conn, err := ln.Accept()
         if err != nil {
               log.Fatal(err)
         }
         go fs.readLoop(conn)
     }
}

func (fs  *FileServer) readLoop(conn net.Conn){
      buf  := new(bytes.Buffer)
      for {
         var size int64
         binary.Read(conn, binary.LittleEndian, &size)
         n, err  := io.CopyN(buf,conn,size)
         if err != nil{
             log.Fatal(err)
         }
         fmt.Println(buf.Bytes())
         fmt.Printf("Received %d bytes over the network\n",n)
      }
}

func sendFile(size int) error {
      file := make([]byte, size)
      _, err := io.ReadFull(rand.Reader, file)
      if err != nil {
          return err
      }
      conn, err := net.Dial("tcp",":3000")
      if err != nil {
          return err
      }
      binary.Write(conn, binary.LittleEndian,int64(size))
      n, err := io.CopyN(conn, bytes.NewReader(file),int64(size))
      if err != nil {
          return err
      }

      fmt.Printf("writtern %d bytes over the network\n",n)
      return nil
}

func main(){
     go func() {
         time.Sleep(4* time.Second)
         sendFile(20000)
     }()
     server := &FileServer{}
     server.start()
}
