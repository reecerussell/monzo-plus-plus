import React from "react";
import {
	Container,
	Divider,
	Grid,
	Header,
	Image,
	List,
	Segment,
} from "semantic-ui-react";

const Footer = () => (
	<Segment
		inverted
		vertical
		style={{ margin: "5em 0em 0em", padding: "5em 0em" }}
	>
		<Container textAlign="center">
			<Grid divided inverted stackable>
				<Grid.Column width={3}>
					<Header inverted as="h4" content="Info" />
					<List link inverted>
						<List.Item as="a" href="/">
							Home
						</List.Item>
					</List>
				</Grid.Column>
				<Grid.Column width={3}>
					<Header inverted as="h4" content="Developer" />
					<List link inverted>
						<List.Item as="a" href="/docs">
							API Docs
						</List.Item>
						<List.Item
							as="a"
							href="https://github.com/reecerussell/monzo-plus-plus"
						>
							Source
						</List.Item>
						<List.Item as="a" href="https://reece-russell.co.uk">
							Reece Russell
						</List.Item>
					</List>
				</Grid.Column>
				<Grid.Column width={10}>
					<Header inverted as="h4" content="About" />
					<p>
						Monzo++ is an experimental project, built to explore and
						make use of webhooks through a custom third party
						system. Although experimental, Monzo++ is completely
						functional and can provide users with additional
						features.
					</p>
				</Grid.Column>
			</Grid>

			<Divider inverted section />
			<Image
				centered
				size="mini"
				src="https://lh3.googleusercontent.com/iDeb12CKMVdgDqBD9yJ9UehaWkKXFdPMtuUA8Jt0sOvxXzOm21qNGbA6D5_gdDZtAk4=w300"
			/>
			<List horizontal inverted divided link size="small">
				<List.Item as="a" href="http://reece-russell.co.uk">
					Designed and developed by{" "}
					<span style={{ textDecoration: "underline" }}>
						Reece Russell
					</span>
				</List.Item>
			</List>
		</Container>
	</Segment>
);

export default Footer;
