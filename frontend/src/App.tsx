import { useState } from 'react'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'

import './App.css'

function App() {
  const [username, setUsername] = useState("");
  const [roomID, setRoomID] = useState("ab12");
  const [contentID, setContentID] = useState("ab2b6c94-47ff-41e5-b36f-f66266af8752");
  const [roomDetails, setRoomDetails] = useState(null);

  const roomServiceAddr = 'http://localhost:8081';
  const videoDownloadAddr = 'http://localhost:8080/video/download';

  const fetchRoomDetails = () => {
    fetch(`${roomServiceAddr}/room/details/${roomID}`, {
      method: 'GET',
      mode: 'cors',
    }).then(resp => {
      if (resp.status == 200) {
        setRoomDetails(resp.json());
      } else {
        console.error('non-200 response to joining room:', resp.status);
      }
    })
  };



  const leaveRoom = () => {
    fetch(`${roomServiceAddr}/room/leave/${username}`, {
      method: 'POST',
      mode: 'cors',
    }).then(resp => {
      if (resp.status == 200) {
        setRoomDetails(null);
      } else {
        console.error('non-200 response to leaving room:', resp.status);
      }
    })
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    fetch(`${roomServiceAddr}/room/join/${roomID}/${username}`, {
      method: 'POST',
      mode: 'cors',
    }).then(resp => {
      if (resp.status == 200) {
        console.info('joined room!')
        fetchRoomDetails()
      }
    })
  };

  return (
    <div className="App">
      <h1>FLMnCHLL</h1>
      <h2>Due to time constraints, this is a PoC.</h2>
      <h2>There is a pre-existing room with short-id: ab12</h2>
      { roomDetails == null ? 
        (
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1">
              <Form.Label>Enter room id</Form.Label>
              <Form.Control 
                as="textarea" 
                rows={1} 
                placeholder='id'
                value={roomID}
                onChange={(e) => setRoomID(e.target.value) }
              />
              <Form.Label>Enter username</Form.Label>
              <Form.Control 
                as="textarea" 
                rows={1} 
                placeholder='username'
                value={username}
                onChange={(e) => setUsername(e.target.value) }
              />
              <Button variant="primary" type="submit">
                Enter the room
              </Button>
            </Form.Group>
          </Form>
        )
      :
        (
          <>
          <p>In room</p>
          <Button variant="primary" onClick={leaveRoom}>
            Leave room
          </Button>

          <div>
            <video controls>
              <source src={`${videoDownloadAddr}/${contentID}?enc=webm`} type="video/webm" />
              <source src={`${videoDownloadAddr}/${contentID}?enc=mp4`} type="video/mp4" />
              <p>
                Your browser doesn't support HTML video. Here is a
                <a href="myVideo.mp4">link to the video</a> instead.
              </p>
            </video>
          </div>
          </>
        )
      }
      
     
    </div>
  )
}

export default App
