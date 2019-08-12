#!/usr/bin/perl

use 5.018;
use Data::Dumper;

my $pid = $ARGV[0];
($pid) = (`ps` =~ m|^(\d+)\s+.+s/$1/insolar|gm) if $pid=~/n(\d+)/i;

my %g;
my $g = "";
my $l = 0;


for (split /\n/, `gops stack $pid`) {
    chomp;
    if (/^goroutine/) {
        $g{$g} += 1;
        $g = "";
        $l = 0;
        next;
    }
    $l++;
    (undef, $g) = split /\s+/ if 2 == $l;
}
say "$g{$_}\t $_" for sort { $g{$b} <=> $g{$a}} keys %g;
