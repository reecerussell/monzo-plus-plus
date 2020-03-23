import React from "react";
import { Redirect } from "react-router-dom";
import { Loader, Message, Form, Button } from "semantic-ui-react";

const Login = ({
	username,
	password,
	handleUpdateUsername,
	handleUpdatePassword,
	handleSubmit,
	loading,
	error,
	redirect,
}) => {
	if (redirect) {
		return <Redirect to={redirect} />;
	}

	return (
		<>
			<Form onSubmit={handleSubmit} error={error !== null}>
				<Loader active={loading} />
				<Message
					error
					content={error}
					header="Oops, something went wrong!"
				/>
				<Form.Field>
					<label htmlFor="username">Username</label>
					<input
						type="text"
						name="username"
						id="username"
						autoComplete="username current-username"
						placeholder="username"
						value={username}
						onChange={handleUpdateUsername}
					/>
				</Form.Field>
				<Form.Field>
					<label htmlFor="password">Password</label>
					<input
						type="password"
						name="password"
						id="password"
						autoComplete="password"
						placeholder="Password..."
						value={password}
						onChange={handleUpdatePassword}
					/>
				</Form.Field>
				<Button type="submit">Login</Button>
			</Form>
		</>
	);
};

export default Login;
