import { useState } from 'react'
import './App.css'

function App() {
  

  return (
    <div className="App">
      <h1>FLMnCHLL</h1>
      <div className="card">
        <h2>Create account</h2>
        <form method="post" action="https://account.flmnchll.lassiheikkila.com/account/signup" enctype="multipart/form-data">
          <div>
            <label for="username">Username: </label>
            <input type="text" id="username" name="username" /><br />

            <label for="password">Password: </label>
            <input type="password" id="password" name="password" /><br />
          </div>

          <div>
            <button>Submit</button>
          </div>
        </form>
      </div>
      <div className="card">
        <h2>Log in</h2>
        <form method="post" action="https://account.flmnchll.lassiheikkila.com/account/signup" enctype="multipart/form-data">
          <div>
            <label for="username">Username: </label>
            <input type="text" id="username" name="username" /><br />

            <label for="password">Password: </label>
            <input type="password" id="password" name="password" /><br />
          </div>
          <div>
            <button>Login</button>
          </div>
        </form>
      </div>
      <div>
        <video controls>
          <source src="http://localhost:8080/video/download/ab2b6c94-47ff-41e5-b36f-f66266af8752?enc=webm" type="video/webm" />
          <source src="http://localhost:8080/video/download/ab2b6c94-47ff-41e5-b36f-f66266af8752?enc=mp4" type="video/mp4" />
          <p>
            Your browser doesn't support HTML video. Here is a
            <a href="myVideo.mp4">link to the video</a> instead.
          </p>
        </video>
      </div>
    </div>
  )
}

export default App
