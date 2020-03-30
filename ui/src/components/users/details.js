import React from "react";
import { Message, Form, Button, ButtonGroup } from "semantic-ui-react";
import DeleteContainer from "../../containers/users/delete";
import PropTypes from "prop-types";

const propTypes = {
	id: PropTypes.string,
	data: PropTypes.object,
	roles: PropTypes.array,
	loading: PropTypes.bool,
	error: PropTypes.string,
	success: PropTypes.string,
	readonly: PropTypes.bool,
	handleSubmit: PropTypes.func,
	handleUpdate: PropTypes.func,
	toggleMode: PropTypes.func,
};
const defaultProps = {
	roles: [],
	error: null,
	success: null,
};

const Details = ({
	id,
	data,
	loading,
	error,
	success,
	readonly,
	handleSubmit,
	handleUpdate,
	toggleMode,
}) => {
	const dateEnabled = (data
	? data.enabled
	: false)
		? new Date(data.dateEnabled)
		: null;

	return (
		<Form
			onSubmit={handleSubmit}
			error={error !== null}
			success={success !== null}
			warning={data ? !data.enabled : false}
		>
			<Message error header="An error occurred!" content={error} />
			<Message success content={success} />
			<Message
				warning
				header="User not enabled!"
				content="This user is either currently disabled or has not yet been enabled."
			/>

			<Form.Field>
				<label htmlFor="username">Username</label>
				<input
					type="text"
					name="username"
					id="username"
					value={data ? data.username : ""}
					onChange={handleUpdate}
					autoComplete="off"
					autoCapitalize="off"
					placeholder="Enter a username..."
					readOnly={readonly}
				/>
			</Form.Field>

			{(data ? (
				data.enabled
			) : (
				false
			)) ? (
				<p>
					This user was enabled on {dateEnabled.toLocaleDateString()}{" "}
					at {dateEnabled.toLocaleTimeString()}
				</p>
			) : null}

			<Form.Field>
				<ButtonGroup>
					{readonly ? (
						<Button
							type="button"
							onClick={toggleMode}
							color="green"
						>
							Edit
						</Button>
					) : (
						<>
							<Button
								type="submit"
								loading={loading}
								color="green"
							>
								Save
							</Button>
							<Button type="button" onClick={toggleMode}>
								Cancel
							</Button>
							<DeleteContainer id={id} />
						</>
					)}
				</ButtonGroup>
			</Form.Field>
		</Form>
	);
};

Details.propTypes = propTypes;
Details.defaultProps = defaultProps;

export default Details;
