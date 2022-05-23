import logo from '../logo.svg';
import './App.css';

import React, { useEffect, useState } from 'react';

import { Login } from '../Login';
import { Profile } from '../Profile';

const LS_KEY = 'login-with-metamask:auth';

export const App = () => {
  const [state, setState] = useState({});

	useEffect(() => {
		// Access token is stored in localstorage
		const ls = window.localStorage.getItem(LS_KEY);
		const auth = ls && JSON.parse(ls);
		setState({ auth });
	}, []);

	const handleLoggedIn = (auth) => {
		localStorage.setItem(LS_KEY, JSON.stringify(auth));
		setState({ auth });
	};

	const handleLoggedOut = () => {
		localStorage.removeItem(LS_KEY);
		setState({ auth: undefined });
	};

	const { auth } = state;

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Login with metamask demo
        </p>
      </header>
      <div className="App-intro">
        {auth ? (<Profile auth={auth} onLoggedOut={handleLoggedOut} />) : (<Login onLoggedIn={handleLoggedIn} />)}
      </div>
    </div>
  );
};
