import React from "react";
import { Loader, Header, Form, Message, Button, Grid } from "semantic-ui-react";
import Layout from "./layout";

export default function ChangePassword({
	loading,
	error,
	success,
	data,
	username,
	handleUpdateCurrentPassword,
	handleUpdateNewPassword,
	handleUpdateConfirmPassword,
	handleSubmit,
}) {
	const { currentPassword, newPassword, confirmPassword } = data;

	const currentPasswordError = currentPassword.error
		? {
				content: currentPassword.error,
				pointing: "above",
		  }
		: null;
	const newPasswordError = newPassword.error
		? {
				content: newPassword.error,
				pointing: "above",
		  }
		: null;
	const confirmPasswordError = confirmPassword.error
		? {
				content: confirmPassword.error,
				pointing: "above",
		  }
		: null;

	return (
		<Layout>
			<Loader active={loading} />

			<Header as="h1">Change Password</Header>

			<Grid stackable>
				<Grid.Row>
					<Grid.Column width={8}>
						<Form
							onSubmit={handleSubmit}
							error={error !== null}
							success={success}
						>
							<Message
								error
								header="An error occured!"
								content={error}
							/>

							<Message
								success
								header="Success!"
								content={success}
							/>

							<Form.Field>
								<input
									type="hidden"
									value={username}
									autoComplete="username current-username"
								/>
							</Form.Field>

							<Form.Field error={currentPasswordError}>
								<label>Current password</label>
								<input
									placeholder="Current password..."
									onChange={handleUpdateCurrentPassword}
									type="password"
									autoComplete="current-password"
								/>
							</Form.Field>
							<Form.Field error={newPasswordError}>
								<label>New password</label>
								<input
									placeholder="Enter a new password..."
									onChange={handleUpdateNewPassword}
									type="password"
									autoComplete="new-password"
								/>
							</Form.Field>
							<Form.Field error={confirmPasswordError}>
								<label>Confirm password</label>
								<input
									placeholder="Confirm password..."
									onChange={handleUpdateConfirmPassword}
									type="password"
								/>
							</Form.Field>

							<Form.Field
								control={Button}
								content="Change Password"
							/>
						</Form>
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Layout>
	);
}
