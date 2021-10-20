import logo from '../assets/logo.svg';
import './App.scss';
import { LoginPage } from './login/LoginPage';
import { Routes, Route } from 'react-router';
import { Link } from 'react-router-dom';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <nav>
          <Link to="/">Home</Link>
          <Link to="/login">Log in</Link>
        </nav>
        <Routes>
          <Route path="login" element={<LoginPage />}></Route>
        </Routes>
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
