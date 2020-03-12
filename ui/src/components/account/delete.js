import React from "react";
import { Form, Button, Loader, Message } from "semantic-ui-react";
import { Redirect } from "react-router-dom";

const Delete = ({ handleSubmit, error, loading, redirect }) => {
	if (redirect) {
		return <Redirect to={redirect} />;
	}

	return (
		<Form onSubmit={handleSubmit} error={error}>
			<Loader active={loading} />
			<Message error header="An error occured!" content={error} />
			<Form.Field>
				<Button color="red" type="submit">
					Delete Account
				</Button>
			</Form.Field>
		</Form>
	);
};

export default Delete;
