import './RegisterPage.scss';
import axios from 'axios';
import { Form, Button } from 'react-bootstrap';
import { useState } from 'react';

export function RegisterPage() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [repeatPassword, setRepeatPassword] = useState('')

    const handleSubmit = async (event) => {
        event.preventDefault();

        if (password !== repeatPassword) {
            alert('Repeat password incorrect!');
            return;
        }

        let enc = new TextEncoder();
        let hash = await crypto.subtle.digest('SHA-256', enc.encode(password));
        
        const packed = window.btoa(
            String.fromCharCode.apply(null, new Uint8Array(hash))
        )

        let { data } = await axios.post('/api/account/register', {
                username,
                password: packed,
            });

        console.log(data);
        if (data.success) {
            alert('Register successful!')
        } else {
            alert('Not registered!');
        }
    };

    return (
        <>
            <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3">
                    <Form.Control 
                        type="text" 
                        placeholder="Enter username"
                        value={username}
                        onChange={e => setUsername(e.currentTarget.value)}
                    />
                </Form.Group>
                <Form.Group className="mb-3">
                    <Form.Control 
                        type="password" 
                        placeholder="Enter password"
                        value={password}
                        onChange={e => setPassword(e.currentTarget.value)}
                        />
                </Form.Group>
                <Form.Group className="mb-3">
                    <Form.Control 
                        type="password" 
                        placeholder="Repeat your password"
                        value={repeatPassword}
                        onChange={e => setRepeatPassword(e.currentTarget.value)}
                        />
                </Form.Group>
                <Button variant="primary" type="submit" size="md">
                    Sign Up
                </Button>
            </Form>
        </>
    )
}