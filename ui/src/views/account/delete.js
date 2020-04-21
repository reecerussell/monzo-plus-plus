import React from "react";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import DeleteContainer from "../../containers/account/delete";
import Layout from "../../components/account/layout";

const Delete = () => (
	<Authorise>
		<Layout>
			<Header as="h2" color="red">
				Delete your account
			</Header>
			<Grid>
				<Grid.Row>
					<Grid.Column width={8}>
						<p>
							This action is permanent and cannot be reversed. All
							plugins and features will no longer be active with
							your account.
						</p>
						<p>
							<b>Are you sure?</b>
						</p>
						<DeleteContainer />
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Layout>
	</Authorise>
);

export default Delete;
