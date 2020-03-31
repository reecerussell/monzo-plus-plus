import React from "react";
import { useParams, Link } from "react-router-dom";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import EditContainer from "../../containers/roles/edit";
import PermissionsContainer from "../../containers/roles/permissions";

const List = () => {
	const { id } = useParams();

	return (
		<Authorise roles={["Admin", "Role Manager"]}>
			<Header as="h1">Edit</Header>

			<Grid stackable>
				<Grid.Row>
					<Grid.Column width="5">
						<p>Use this form to manage this role.</p>

						<p>
							<Link to="/roles">Back to roles</Link>
						</p>

						<EditContainer id={id} />
					</Grid.Column>
					<Grid.Column width="11">
						<PermissionsContainer id={id} />
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Authorise>
	);
};

export default List;
