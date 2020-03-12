import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import { useParams, Link } from "react-router-dom";

import DetailsContainer from "../../containers/users/details";

const Users = () => {
	const { id } = useParams();

	return (
		<Authorise roles={["Admin"]}>
			<Header as="h1">User Details</Header>
			<p>
				<Link to="/users">View all users</Link>
			</p>
			<DetailsContainer id={id} />
		</Authorise>
	);
};

export default Users;
