import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import PendingContainer from "../../containers/users/pending";
import { Link } from "react-router-dom";

const Pending = () => (
	<Authorise roles={["Admin"]}>
		<Header as="h1">Pending Users</Header>

		<p>
			<Link to="/users">View all users</Link>
		</p>

		<PendingContainer />
	</Authorise>
);

export default Pending;
