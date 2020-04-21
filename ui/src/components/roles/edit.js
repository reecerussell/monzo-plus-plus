import React from "react";
import { Form, Message, Button } from "semantic-ui-react";
import PropTypes from "prop-types";
import DeleteContainer from "../../containers/roles/delete";

const propTypes = {
	id: PropTypes.string.isRequired,
	data: PropTypes.object.isRequired,
	error: PropTypes.string,
	success: PropTypes.bool,
	loading: PropTypes.bool,
	handleSubmit: PropTypes.func.isRequired,
	handleUpdate: PropTypes.func.isRequired,
	updateError: PropTypes.func.isRequired,
};
const defaultProps = {
	error: null,
	success: false,
	loading: false,
};

const Edit = ({
	id,
	data,
	handleUpdate,
	handleSubmit,
	error,
	success,
	loading,
	updateError,
}) => {
	return (
		<Form onSubmit={handleSubmit} error={error !== null} success={success}>
			<Message error header="An error occured!" content={error} />
			<Message
				success
				header="Saved!"
				content="Role saved successfully!"
			/>

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

			<Form.Field>
				<Button loading={loading} color="green">
					Save
				</Button>
				<DeleteContainer id={id} onError={updateError} />
			</Form.Field>
		</Form>
	);
};

Edit.propTypes = propTypes;
Edit.defaultProps = defaultProps;

export default Edit;
