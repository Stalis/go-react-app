import './LoginPage.scss';
import { Form, Button } from 'react-bootstrap';
import React from 'react';
import axios from 'axios';

export class LoginPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            username: '',
            password: '',
        };
    }

    handleSubmit = (event) => {
        event.preventDefault();
        axios.post('/api/login', {
            username: this.state.username,
            password: this.state.password,
        }).then(response => console.log(response.data));
    };

    inputChangeHandler = (name) => (event) => {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;

        this.setState({
            [name]: value,
        });
    }

    render() {
        return (
            <Form onSubmit={this.handleSubmit}>
                <Form.Group className="mb-3" controlId="formBasicUsername">
                    <Form.Label>Username</Form.Label>
                    <Form.Control type="text" placeholder="Enter username" 
                        name="username"
                        value={this.state.username}
                        onChange={this.inputChangeHandler('username')} />
                    <Form.Text className="text-muted">
                        We'll never share your username with anyone else.
                    </Form.Text>
                </Form.Group>

                <Form.Group className="mb-3" controlId="formBasicPassword">
                    <Form.Label>Password</Form.Label>
                    <Form.Control type="password" placeholder="Password"
                        name="password"
                        value={this.state.password}
                        onChange={this.inputChangeHandler('password')} />
                </Form.Group>
                <Form.Group className="mb-3" controlId="formBasicCheckbox">
                    <Form.Check type="checkbox" label="Check me out" />
                </Form.Group>
                <Button variant="primary" type="submit">
                    Log In
                </Button>
            </Form>
        );
    }
}