import React from "react";
import LoginContainer from "../../containers/login/login";
import { Grid, Header } from "semantic-ui-react";
import { Link } from "react-router-dom";

const Login = () => (
	<>
		<Header as="h1">Login</Header>
		<p>Login to access your plugins and account.</p>
		<p>
			<Link to="/register">Don't have an account?</Link>
		</p>

		<Grid>
			<Grid.Row>
				<Grid.Column width={5}>
					<LoginContainer />
				</Grid.Column>
			</Grid.Row>
		</Grid>
	</>
);

export default Login;
