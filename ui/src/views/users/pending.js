import React from "react";
import { Header } from "semantic-ui-react";
import PendingContainer from "../../containers/users/pending";
import { Link } from "react-router-dom";

const Pending = () => (
	<>
		<Header as="h1">Pending Users</Header>

		<p>
			<Link to="/users">View all users</Link>
		</p>

		<PendingContainer />
	</>
);

export default Pending;
