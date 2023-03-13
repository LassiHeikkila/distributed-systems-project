import { useState, useEffect, useCallback } from 'react';
import { BrowserRouter, Route, Switch, useHistory } from 'react-router-dom';

import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import Stack from 'react-bootstrap/Stack';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';

import PeerJs from 'peerjs';

import './App.css'

let peerClient: PeerJs;
let connection: PeerJs.DataConnection;

let g_username:string;
let g_roomID:string;

const roomServiceAddr = 'http://localhost:8081';
const videoDownloadAddr = 'http://localhost:8080/video/download';
const videoSearchAddr = 'http://localhost:8080/video/search'

interface ChatMessage {
  id: number;
  sender: string;
  message: string;
  time: string;
};

interface ControlCommand {
  id: number;
  time: string;
  play: boolean;
  pause: boolean;
};

function LandingView() {
  const [username, setUsername] = useState(g_username);
  const [roomID, setRoomID] = useState(g_roomID);
  const history = useHistory();

  const handleSubmit = (event: Event) => {
    event.preventDefault();

    fetch(`${roomServiceAddr}/room/join/${roomID}/${username}`, {
      method: 'POST',
      mode: 'cors',
    }).then(resp => {
      if (resp.status == 200) {
        console.info('joined room!');
        history.replace("/room");
      }
    })
  };

  useEffect(() => {
    g_username = username;
  }, [username]);

  useEffect(() => {
    g_roomID = roomID;
  }, [roomID]);
  

  return (
    <>
    <Container fluid='lg'>
      <Row>
        <Col>
        <h3>There is a pre-existing room with short-id: ab12</h3>
        <h4>Use the form to join the room:</h4>
        <Form onSubmit={handleSubmit}>
          <Form.Group size='lg' className="mb-3">
            <Form.Label>Set room id</Form.Label>
            <Form.Control 
              type="text"
              placeholder='id'
              value={roomID}
              onChange={(e) => setRoomID(e.target.value) }
            />
            <br />
            <Form.Label>Set username</Form.Label>
            <Form.Control 
              type="text"
              placeholder='username'
              value={username}
              onChange={(e) => setUsername(e.target.value) }
            />
            <br />
            <Button variant="primary" type="submit">
              Enter the room
            </Button>
          </Form.Group>
        </Form>
        </Col>
      </Row>
    </Container>
    </>
  )
}

function VideoPlayer(props:{selectedContentID:string}) {
  return (
    (props.selectedContentID != "") ? (
      <video controls width="800px" >
      <source src={`${videoDownloadAddr}/${props.selectedContentID}?enc=webm`} type="video/webm" />
      <source src={`${videoDownloadAddr}/${props.selectedContentID}?enc=mp4`} type="video/mp4" />
      <p>
        Your browser doesn't support HTML video. Here is a
        <a href="myVideo.mp4">link to the video</a> instead.
      </p>
    </video>
    ) : (<p>Select content to start watching</p>)
  )
}

function RoomView() {
  const history = useHistory();
  const [roomID, setRoomID] = useState(g_roomID);
  const [username, setUsername] = useState(g_username);

  const [peerServer, setPeerServer] = useState("");
  const [peerServerConnected, setPeerServerConnected] = useState(false);
  const [availablePeerClient, setAvailablePeerClient] = useState(peerClient);

  const [peerConnections, setPeerConnections] = useState<PeerJs.DataConnection>([]);
  
  const [roomDetailsTimer, setRoomDetailsTimer] = useState({});
  const [availableContentTimer, setAvailableContentTimer] = useState({});

  const [availableContent, setAvailableContent] = useState([]);
  const [selectedContentID, setSelectedContentID] = useState("");

  const [roomDetails, setRoomDetails] = useState({});

  const [messages, setMessages] = useState([]);

  const appendMessage = useCallback((msg:ChatMessage) => {
    setMessages((msgs) => [...msgs, msg]);
  }, []);
  
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
        console.info("successfully left room");
    })
    .catch(err => {
        console.error('error leaving room:', err);
    })

    setRoomDetails({});
    clearInterval(roomDetailsTimer);
    clearInterval(availableContentTimer);
    history.replace("/");
  };

  const sendMessage = (msg:string) => {
    const m:ChatMessage = {
      id: Date.now(),
      message: msg,
      sender: username,
      time: new Date().toISOString(),
    }
    appendMessage(m);

    for (const peer in peerConnections) {
      // peer.send(m);
    }
  };

  useEffect(() => {
    if (roomDetails.peerServerAddr != "") {
      setPeerServer(roomDetails.peerServerAddr);
    }
  }, [roomDetails]);

  useEffect(() => {
    if (peerServer != "" && !peerServerConnected) {
      setAvailablePeerClient(new PeerJs(username, {
        host: peerServer,
        path: '/',
        port: 9000,
      }))
    }
  }, [username, peerServer, peerServerConnected]);

  useEffect(() => {
    fetchRoomDetails();
    fetchAvailableContent();

    const roomDetailsTimerID = setInterval(fetchRoomDetails, 10000);
    setRoomDetailsTimer(roomDetailsTimerID);

    const availableContentTimerID = setInterval(fetchAvailableContent, 30000);
    setAvailableContentTimer(availableContentTimerID);
  }, []);

  const submitMessage = (evt:React.FormEventHandler<HTMLFormElement>) => {
    const msg = evt.currentTarget.elements.namedItem('message_field').value;
    sendMessage(msg);
  };
 
  return (
    <Container>
      <Row>
        <Col>
          <h3><>In room "{roomID}" as user "{username}"</></h3>
        </Col>
        <Col>
          <Button variant="primary" onClick={leaveRoom}>
            Leave room
          </Button>
        </Col>
      </Row>
      <Row>
        <VideoPlayer selectedContentID={selectedContentID} />
      </Row>
      <Row>
      Users present: 
      { roomDetails.users ? roomDetails.users.map((user) => (
        <> {user.name} </>
      )) : <></> }
      </Row>
      <Row>
      <p>Available content:
      {
        availableContent.length > 0 ? availableContent.map((content) => (
          <Button onClick={() => {setSelectedContentID(content.contentID)}} >
            {content.name}
          </Button>
        )) : <></> }
      </p>
      </Row>
      <Row>
        Chat:
        <Container>
          {messages.map((msg:ChatMessage) => (
            <p key={msg.id} style={{ color: msg.sender == username ? '#999' : '#222' }}>
              <b>{msg.sender}</b> ({msg.time}): {msg.message}
            </p>
          ))}
        </Container>
        <Form onSubmit={submitMessage}>
          <Form.Control id="message_field" />
          <Form.Text muted>New message</Form.Text>
        </Form>
      </Row>
    </Container>
  )
}

function App() {
  return (
    <Container fluid>
      <Navbar fixed="top">
        <Container>
          <Navbar.Brand href="https://github.com/LassiHeikkila/flmnchll">
            <b>FLMnCHLL</b>
          </Navbar.Brand>
          <Navbar.Text><i>Due to time contraints, this is a simple PoC</i></Navbar.Text>
        </Container>
      </Navbar>
      <BrowserRouter>
        <Switch>
          <Route exact path="/" component={LandingView} />
          <Route exact path="/room" component={RoomView} />
        </Switch>
      </BrowserRouter>
    </Container>
  );
};

export default App;
