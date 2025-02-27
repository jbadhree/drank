import { NextResponse } from 'next/server';

export async function GET() {
  return NextResponse.json({
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'Not set',
    // Don't include sensitive environment variables here
  });
}
