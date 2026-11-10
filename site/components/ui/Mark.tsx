const basePath = process.env.NEXT_PUBLIC_BASE_PATH || "";

export default function Mark({ className = "" }: { className?: string }) {
  return (
    <img
      src={`${basePath}/logo.png`}
      alt=""
      aria-hidden="true"
      width={28}
      height={28}
      className={`rounded-[6px] object-cover ${className}`}
    />
  );
}
