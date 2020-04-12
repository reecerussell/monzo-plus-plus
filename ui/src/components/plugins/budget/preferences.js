import React, { useRef } from "react";
import {
	Modal,
	Form,
	Message,
	Button,
	Icon,
	ButtonGroup,
} from "semantic-ui-react";

const Preferences = ({
	monthlyBudget,
	handleUpdateBudget,
	error,
	success,
	loading,
	handleUpdate,
	showModal,
	toggleModal,
}) => {
	const formRef = useRef(null);

	return (
		<Modal
			trigger={
				<Button color="blue" type="button" onClick={toggleModal}>
					<Icon name="pencil" />
				</Button>
			}
			open={showModal}
			onClose={toggleModal}
		>
			<Modal.Header>Preferences</Modal.Header>

			<Modal.Content>
				<Form
					ref={formRef}
					onSubmit={handleUpdate}
					error={error !== null}
					success={success !== null}
				>
					<Message error header="An error occured!" content={error} />
					<Message success content={success} />

					<p>
						Use this form to manage your plugin preferences. By
						changing your monthly budget below will affect your
						daily budget amount.
					</p>

					<Form.Field>
						<label htmlFor="monthlyBudget">Monthly Budget</label>
						<input
							id="monthlyBudget"
							name="monthlyBudget"
							type="number"
							min="0"
							step="0.01"
							value={monthlyBudget}
							onChange={handleUpdateBudget}
							placeholder="300.00"
						/>
					</Form.Field>
				</Form>
			</Modal.Content>

			<Modal.Actions>
				<ButtonGroup>
					<Button color="grey" type="button" onClick={toggleModal}>
						Close
					</Button>
					<Button
						color="green"
						type="button"
						onClick={() => {
							if (formRef.current) {
								formRef.current.handleSubmit();
							}
						}}
						loading={loading}
					>
						Save
					</Button>
				</ButtonGroup>
			</Modal.Actions>
		</Modal>
	);
};

export default Preferences;
