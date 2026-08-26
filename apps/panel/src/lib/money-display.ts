import { Moneda } from "@/providers/graphql/graphql";


const LOCALES: Partial<Record<Moneda, { locale: string; currency: string }>> = {
	[Moneda.Usd]: { locale: 'en-US', currency: 'USD' }
}

export function money(value: number, currency: Moneda = Moneda.Usd): string {
	if (currency === Moneda.Ved) {
		const formatted = new Intl.NumberFormat('es-VE', {
			minimumFractionDigits: 2,
			maximumFractionDigits: 2
		}).format(value)
		return `Bs. ${formatted}`
	}

	const config = LOCALES[currency] ?? LOCALES[Moneda.Usd]!
	return new Intl.NumberFormat(config.locale, {
		style: 'currency',
		currency: config.currency
	}).format(value)
}