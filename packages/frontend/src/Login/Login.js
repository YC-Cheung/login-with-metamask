import './Login.css'

import React, { useState } from 'react';
import Web3 from 'web3';

let web3 = Web3 | undefined


export const Login = ({ onLoggedIn }) => {
  const [loading, setLoading] = useState(false);

  const handleAuthenticate = (data) => {
    let { publicAddress, signature, nonce } = data
    fetch(`${process.env.REACT_APP_BACKEND_URL}/user/auth`, {
      body: JSON.stringify({publicAddress, signature, nonce}),
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
    }).then(response => response.json())
    .then((data) => {
      if (data.accessToken) {
        setLoading(false);
        onLoggedIn(data);
      }
    });
  };

  const handleSignMessage = async (publicAddress, nonce) => {
    try {
      const signature = await web3.eth.personal.sign(
        `I am signing my one-time nonce: ${nonce}`,
        publicAddress,
        ''
      );

      return { publicAddress, signature, nonce };
    } catch (err) {
      throw new Error(
        'You need to sign the message for log in.'
      );
    }
  };

  const handleClick = async () => {
    if (typeof window.ethereum == 'undefined') {
      window.alert('Please install MetaMask first.');
      return;
    }

    if (!web3) {
      try {
        await window.ethereum.enable();
        web3 = new Web3(window.ethereum);
      } catch (error) {
        window.alert('You need to allow MetaMask.');
        return;
      }
    }

    const coinbase = await web3.eth.getCoinbase();
    if (!coinbase) {
      window.alert('Please activate MetaMask first');
      return;
    }

    const publicAddress = coinbase.toLowerCase();
    setLoading(true);

    fetch(
			`${process.env.REACT_APP_BACKEND_URL}/user/nonce?publicAddress=${publicAddress}`
		)
			.then((response) => response.json())
			.then((res) =>
        handleSignMessage(publicAddress, res.nonce)
			)
			.then(handleAuthenticate)
      .then(onLoggedIn)
			.catch((err) => {
				window.alert(err);
				setLoading(false);
			});
  }

  return (
    <div>
      <button className="login-btn login-mm" onClick={handleClick}>
        {loading ? 'Loading...' : 'Login with MetaMask'}
      </button>
    </div>
  );
};