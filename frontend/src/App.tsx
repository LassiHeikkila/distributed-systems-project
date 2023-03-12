import { useState, useEffect } from 'react'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'

import './App.css'

function App() {
  const [username, setUsername] = useState("");
  const [roomID, setRoomID] = useState("");
  const [contentID, setContentID] = useState("");
  const [availableContent, setAvailableContent] = useState([]);
  const [roomDetails, setRoomDetails] = useState(null);
  const [peerServer, setPeerServer] = useState("");

  const [timer, setTimer] = useState(0);

  const roomServiceAddr = 'http://localhost:8081';
  const videoDownloadAddr = 'http://localhost:8080/video/download';
  const videoSearchAddr = 'http://localhost:8080/video/search'

  const fetchRoomDetails = () => {
    fetch(`${roomServiceAddr}/room/details/${roomID}`, {
      method: 'GET',
      mode: 'cors',
    })
    .then(resp => resp.json())
    .then(data => {
      setRoomDetails(data);
    })
    .catch(err => {
      console.error("error fetching room details:", err)
    })
  };

  const fetchAvailableContent = () => {
    fetch(`${videoSearchAddr}`)
      .then(resp => resp.json())
      .then(data => {
        setAvailableContent(data);
      })
      .catch(err => {
        console.error('error searching content:', err)
      })
  }

  const leaveRoom = () => {
    fetch(`${roomServiceAddr}/room/leave/${username}`, {
      method: 'POST',
      mode: 'cors',
    })
    .then(resp => resp.json())
    .then(data => {
        setRoomDetails(null);
        setRoomID("");
        clearInterval(timer);
    })
    .catch(err => {
        console.error('error leaving room:', err);
    })
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    fetch(`${roomServiceAddr}/room/join/${roomID}/${username}`, {
      method: 'POST',
      mode: 'cors',
    }).then(resp => {
      if (resp.status == 200) {
        console.info('joined room!');
        fetchRoomDetails();
        let timerId = setInterval(fetchRoomDetails, 5000);
        setTimer(timerId);
      }
    })
  };

  useEffect(() => {
    if (availableContent.length > 0) {
      setContentID(availableContent[0]['contentID']);
    }
  }, [availableContent]);

  useEffect(() => {
    if (roomDetails == null) {
      setPeerServer("");
    } else {
      setPeerServer(roomDetails.peerServerAddr);
    }
  }, [roomDetails]);

  useEffect(() => {
    fetchAvailableContent();
  }, []);

  return (
    <div className="App">
      <h1>FLMnCHLL</h1>
      <h2>Due to time constraints, this is a very simple PoC.</h2>
      <h2>There is a pre-existing room with short-id: ab12</h2>
      { roomDetails == null ? 
        (
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1">
              <Form.Label>Enter room id</Form.Label>
              <Form.Control 
                type="text"
                placeholder='id'
                value={roomID}
                onChange={(e) => setRoomID(e.target.value) }
              />
              <Form.Label>Enter username</Form.Label>
              <Form.Control 
                type="text"
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
          <p>Users present: 
          { roomDetails.users.map((user) => (
            <> {user.name} </>
          )) }
          </p>
          </>
        )
      }
      
     
    </div>
  )
}

export default App
