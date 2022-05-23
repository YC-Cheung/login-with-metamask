import jwtDecode from 'jwt-decode';
import React, { useState, useEffect } from 'react';
import Blockies from 'react-blockies';

export const Profile = ({ auth, onLoggedOut }) => {
    const [state, setState] = useState({
		loading: false,
        user: {
            uid: undefined,
            publicAddress: '',
            username: '',
        }
	});

	useEffect(() => {
		const { accessToken } = auth;

		fetch(`${process.env.REACT_APP_BACKEND_URL}/user/profile`, {
			headers: {
				Authorization: `Bearer ${accessToken}`,
			},
		})
			.then((response) => response.json())
			.then((data) => {
                let { code, data: user, msg } = data;
                
                if (code === 0) {
                    setState({ ...state, user });
                    return;    
                }

                if (code === 4001) {
                    window.alert(msg);
                    onLoggedOut();
                    return;
                }
            })
			.catch(window.alert);
	}, []);

    const handleChange = ({
		target: { value }
	}) => {
		setState({ ...state, username: value });
	};

    const handleSubmit = () => {
		const { accessToken } = auth;
		const { user, username } = state;

		setState({ ...state, loading: true });

		if (!user) {
			window.alert(
				'The user id has not been fetched yet. Please try again in 5 seconds.'
			);
			return;
		}

		fetch(`${process.env.REACT_APP_BACKEND_URL}/user/modify-name`, {
			body: JSON.stringify({ username }),
			headers: {
				Authorization: `Bearer ${accessToken}`,
				'Content-Type': 'application/json',
			},
			method: 'PATCH',
		})
			.then((response) => response.json())
			.then((data) => {
                let { code, data: user, msg } = data;
                if (code === 0) {
                    setState({ ...state, loading: false, user });
                    return;    
                }

                if (code === 4001) {
                    window.alert(msg);
                    onLoggedOut();
                    return;
                }

                window.alert(msg);
                setState({ ...state, loading: false });
            })
			.catch((err) => {
                console.log('error')
				window.alert(err);
				setState({ ...state, loading: false });
			});
	};

    const { accessToken } = auth;

    const { publicAddress } = jwtDecode(accessToken);
    const { loading, user } = state;
    const username = user && user.username;

    return (
		<div className="Profile">
			<p>
				Logged in as <Blockies seed={publicAddress} />
			</p>
			<div>
				My username is {username ? <pre>{username}</pre> : 'not set.'}{' '}
				My publicAddress is <pre>{publicAddress}</pre>
			</div>
			<div>
				<label htmlFor="username">Change username: </label>
				<input name="username" onChange={handleChange} />
				<button disabled={loading} onClick={handleSubmit}>
					Submit
				</button>
			</div>
			<p>
				<button onClick={onLoggedOut}>Logout</button>
			</p>
		</div>
	);
}