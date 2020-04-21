import React from "react";
import { Form, Button, Message, Grid } from "semantic-ui-react";

const getAccountName = (name) => {
	if (name.includes("user_")) {
		return "Current Account";
	}

	return name;
};

const SetAccount = ({
	redirect,
	error,
	loading,
	success,
	selectedAccount,
	accounts,
	handleSubmit,
	handleUpdateAccount,
}) => (
	<Grid stackable>
		<Grid.Column>
			<Form
				onSubmit={handleSubmit}
				error={error !== null}
				success={success !== null}
			>
				<Message error header="An error occured!" content={error} />
				<Message success content={success} />

				<Form.Field>
					<label htmlFor="account">Select an account</label>
					<select
						id="account"
						name="account"
						value={selectedAccount}
						onChange={handleUpdateAccount}
					>
						{accounts.map((acc, idx) => (
							<option value={acc.id} key={idx}>
								{getAccountName(acc.description)}
							</option>
						))}
					</select>
				</Form.Field>

				<Form.Field>
					<Button color="green" type="submit" loading={loading}>
						Save
					</Button>
				</Form.Field>
			</Form>
		</Grid.Column>
	</Grid>
);

export default SetAccount;
