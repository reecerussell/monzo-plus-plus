import React from "react";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import Layout from "../../components/account/layout";
import SetAccountContainer from "../../containers/account/setAccount";

const SetAccount = () => {
	return (
		<Authorise>
			<Layout>
				<Header as="h2">Select an account</Header>

				<Grid stackable>
					<Grid.Column width="8">
						<p>
							To complete your registration and enable Monzo++,
							you need to select one of your Monzo accounts in
							which to use.
						</p>

						<SetAccountContainer />
					</Grid.Column>
				</Grid>
			</Layout>
		</Authorise>
	);
};

export default SetAccount;
