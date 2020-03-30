import React from "react";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import { useParams, Link } from "react-router-dom";

import DetailsContainer from "../../containers/users/details";
import RolesContainer from "../../containers/users/roles";

const Users = () => {
	const { id } = useParams();

	return (
		<Authorise roles={["Admin"]}>
			<Header as="h1">User Details</Header>

			<Grid stackable>
				<Grid.Row>
					<Grid.Column width={5}>
						<p>
							<Link to="/users">View all users</Link>
						</p>
						<DetailsContainer id={id} />
					</Grid.Column>
					<Grid.Column width={11}>
						<RolesContainer id={id} />
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Authorise>
	);
};

export default Users;
