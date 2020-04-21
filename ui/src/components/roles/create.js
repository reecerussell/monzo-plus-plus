import React from "react";
import { Form, Loader, Message, Button } from "semantic-ui-react";
import { Redirect } from "react-router-dom";
import PropTypes from "prop-types";

const propTypes = {
	data: PropTypes.object.isRequired,
	error: PropTypes.string,
	loading: PropTypes.bool,
	redirect: PropTypes.string,
	handleSubmit: PropTypes.func.isRequired,
	handleUpdate: PropTypes.func.isRequired,
};
const defaultProps = {
	error: null,
	loading: false,
	redirect: null,
};

const Create = ({
	data,
	handleUpdate,
	handleSubmit,
	error,
	loading,
	redirect,
}) => {
	if (redirect) {
		return <Redirect to={redirect} />;
	}

	return (
		<Form onSubmit={handleSubmit} error={error !== null}>
			<Loader active={loading} />
			<Message error header="An error occured!" content={error} />

			<Form.Field>
				<label htmlFor="name">Name</label>
				<input
					type="text"
					value={data.name}
					name="name"
					id="name"
					onChange={handleUpdate}
					autoComplete="off"
					placeholder="Enter a name..."
				/>
			</Form.Field>

			<Form.Field control={Button} content="Create" />
		</Form>
	);
};

Create.propTypes = propTypes;
Create.defaultProps = defaultProps;

export default Create;
