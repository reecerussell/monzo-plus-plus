import React from "react";
import { Link } from "react-router-dom";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import CreateContainer from "../../containers/roles/create";

const List = () => (
	<Authorise roles={["Admin"]}>
		<Header as="h1">Create</Header>

		<p>Use this form to create a new role.</p>

		<p>
			<Link to="/roles">Back to roles</Link>
		</p>

		<Grid stackable>
			<Grid.Row>
				<Grid.Column width="5">
					<CreateContainer />
				</Grid.Column>
			</Grid.Row>
		</Grid>
	</Authorise>
);

export default List;
