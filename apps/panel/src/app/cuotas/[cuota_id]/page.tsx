export default async function CuotaPage({
  params,
}: {
  params: Promise<{ cuota_id: string }>;
}) {
  const { cuota_id } = await params;
  return <div>My Post: {cuota_id}</div>;
}
