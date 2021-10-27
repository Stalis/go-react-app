import './LoginPage.scss';
import { Form, Button } from 'react-bootstrap';
import React, { useState } from 'react';
import axios from 'axios';

export function LoginPage() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = async (event) => {
        event.preventDefault();

        const enc = new TextEncoder();
        const hash = await crypto.subtle.digest('SHA-256', enc.encode(password))

        const packed = window.btoa(
            String.fromCharCode.apply(null, new Uint8Array(hash))
        )

        axios.post('/api/account/login', {
                username,
                password: packed,
            }, {
                headers: {
                    'X-Session-Token': 'initial',
                },
            })
            .then(({ data }) => {
                console.log(data);
                localStorage.setItem("session_token", data.sessionToken);
                alert('Login successful!');
            })
            .catch(({ response }) => {
                console.log(response);
                alert('Login error!');
            });
    };

    return (
        <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3" controlId="formBasicUsername">
                <Form.Label>Username</Form.Label>
                <Form.Control 
                    type="text" 
                    placeholder="Enter username" 
                    name="username"
                    value={username}
                    onChange={e => setUsername(e.currentTarget.value)} />
                <Form.Text className="text-muted">
                    We'll never share your username with anyone else.
                </Form.Text>
            </Form.Group>

            <Form.Group className="mb-3" controlId="formBasicPassword">
                <Form.Label>Password</Form.Label>
                <Form.Control 
                    type="password" 
                    placeholder="Password"
                    name="password"
                    value={password}
                    onChange={e => setPassword(e.currentTarget.value)} />
            </Form.Group>

            <Form.Group className="mb-3" controlId="formBasicCheckbox">
                <Form.Check 
                    type="checkbox" 
                    label="Check me out" />
            </Form.Group>
            
            <Button variant="primary" type="submit">
                Log In
            </Button>
        </Form>
    );
}