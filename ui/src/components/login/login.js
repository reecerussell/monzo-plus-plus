import React from "react";
import { Header, Loader, Message, Form, Button, Grid } from "semantic-ui-react";

const Login = ({
	username,
	password,
	handleUpdateUsername,
	handleUpdatePassword,
	handleSubmit,
	loading,
	error,
}) => {
	return (
		<>
			<Header as="h1">Login</Header>
			<p>Use this form to log into your account.</p>

			<Grid>
				<Grid.Row>
					<Grid.Column width={5}>
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
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</>
	);
};

export default Login;
