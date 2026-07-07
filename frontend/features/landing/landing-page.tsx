import { Architecture } from "./architecture";
import { Features } from "./features";
import { Footer } from "./footer";
import { Hero } from "./hero";
import { Hipaa } from "./hipaa";
import { Navbar } from "./navbar";
import { TrustedBy } from "./trusted-by";

export function LandingPage() {
	return (
		<div className="min-h-screen bg-slate-950 text-slate-50">
			<Navbar />
			<Hero />
			<Features />
			<Hipaa />
			<TrustedBy />
			<Architecture />
			<Footer />
		</div>
	);
}
