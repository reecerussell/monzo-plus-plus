import React from "react";
import { Redirect, Link } from "react-router-dom";
import { Loader, Message, Form, Button } from "semantic-ui-react";

const Register = ({
	username,
	password,
	confirm,
	handleUpdateUsername,
	handleUpdatePassword,
	handleUpdateConfirm,
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
			<Form onSubmit={handleSubmit} error={error}>
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
						autoComplete="off"
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
						autoComplete="off"
						placeholder="Password..."
						value={password}
						onChange={handleUpdatePassword}
					/>
				</Form.Field>
				<Form.Field>
					<label htmlFor="confirm">Confirm Password</label>
					<input
						type="password"
						name="confirm"
						id="confirm"
						autoComplete="off"
						placeholder="Confirm your password..."
						value={confirm}
						onChange={handleUpdateConfirm}
					/>
				</Form.Field>
				<Button type="submit">Register</Button>
			</Form>
		</>
	);
};

export default Register;
