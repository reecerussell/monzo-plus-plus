import React from "react";
import { Header, Loader, Message, List, Grid, Button } from "semantic-ui-react";
import Layout from "./layout";
import * as User from "../../utils/user";
import { BaseUrl } from "../../utils/fetch";

export default function Index({ data, error, loading }) {
	if (data === null) {
		return null;
	}

	const dateEnabled =
		data && data.enabled ? new Date(data.dateEnabled) : null;

	const content = error ? (
		<Message error header="An error occured!" content={error} />
	) : data ? (
		<List>
			<List.Item icon="user" content={data.username} />
			<List.Item
				icon={data.enabled ? "check" : "close"}
				content={
					data.enabled
						? `Your account was enabled on ${dateEnabled.toLocaleDateString()}, at ${dateEnabled.toLocaleTimeString()}`
						: "Your account has not yet been enabled."
				}
			/>
		</List>
	) : null;

	return (
		<Layout>
			{data.monzoLinked ? null : (
				<>
					<Message>
						<h3>Your account is not linked to Monzo.</h3>
						<p>
							For plugins to work, you need to link your Monzo
							account.
						</p>
						<p>
							<Button
								as="a"
								href={
									BaseUrl +
									"api/auth/monzo/login?id=" +
									User.GetId()
								}
								style={{
									backgroundColor: "#14233c",
									color: "#fff",
								}}
							>
								Link Monzo
							</Button>
						</p>
					</Message>
				</>
			)}
			<Loader active={loading} />
			<Header as="h2">Account</Header>
			<Grid stackable>
				<Grid.Row>
					<Grid.Column width={12}>
						{content}
						<p>
							Here is some information about your account and the
							status of it. You can manage plugins associated with
							it and manage your account as a whole.
						</p>
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Layout>
	);
}
