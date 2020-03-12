import React from "react";
import { Header, Grid } from "semantic-ui-react";
import { Link } from "react-router-dom";
import RegisterContainer from "../../containers/login/register";

const Register = () => (
	<>
		<Header as="h1">Register</Header>
		<p>Sign up to make use of extra custom features for Monzo.</p>
		<p>
			<Link to="/login">Already have an account?</Link>
		</p>

		<Grid>
			<Grid.Row>
				<Grid.Column width={5}>
					<RegisterContainer />
				</Grid.Column>
			</Grid.Row>
		</Grid>
	</>
);

export default Register;
