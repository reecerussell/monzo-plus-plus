import React from "react";
import { Button, Card, Message, Segment, Loader } from "semantic-ui-react";

const Pending = ({ users, loading, error, handleDelete, handleEnable }) => {
	const errorMessage = error ? (
		<Message error header="An error occured!" content={error} />
	) : null;

	const content =
		users && users.length > 0 ? (
			<Card.Group>
				{users.map((user, key) => (
					<Card key={key}>
						<Card.Content>
							<Card.Header>{user.username}</Card.Header>
							<Card.Meta>Newly registered user</Card.Meta>
							<Card.Description>
								To activate this user, click 'Approve' below,
								otherwise click 'Decline'.
							</Card.Description>
						</Card.Content>
						<Card.Content extra>
							<div className="ui two buttons">
								<Button
									basic
									color="green"
									onClick={(e) => handleEnable(user.id, e)}
								>
									Approve
								</Button>
								<Button
									basic
									color="red"
									onClick={(e) => handleDelete(user.id, e)}
								>
									Decline
								</Button>
							</div>
						</Card.Content>
					</Card>
				))}
			</Card.Group>
		) : (
			<p>There are currently no new users.</p>
		);

	return (
		<>
			{errorMessage}

			<Segment basic>
				<Loader active={loading} />

				{content}
			</Segment>
		</>
	);
};

export default Pending;
