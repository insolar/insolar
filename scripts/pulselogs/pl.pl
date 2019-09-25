#!/bin/env perl

# God thanks this is not rust

=HEAD

    Script for analysing insolar node logs for consistence

=cut

use 5.018;

use Data::Dumper;
use FindBin;
use JSON::XS;


our @FILES = map {"$FindBin::Bin/../../.artifacts/launchnet/logs/discoverynodes/$_/output.log"} (1..5);

my @F = map { LogFile->new($_)->init } @FILES;

while (my $obj = $F[0]->read_line) {
    say( encode_json $obj );
}




package ConsensusAnalyzer;

use 5.018;




package LogFile;

use 5.018;
use JSON::XS;

sub new {
    my $class = shift;
    return bless {
        fname => $_[0],
    };
}

sub init {
    my $self = shift;
    my $fn = $self->{fname};
    open $self->{fd}, '<', $fn or die "Can't open $fn : $!";
    $self->read_head;
    return $self;
}

sub DESTROY {
    my $self = shift;
    close $self->{fd};
}

sub read_head {
    my $self = shift;
    my $fd = $self->{fd};
    while(<$fd>) {
        $self->{headbuff} = decode_json($_) , last if /^\{/;
        $self->{head} .= $_;
    }
}

sub read_line {
    my $self = shift;
    if (exists $self->{headbuff}) {
        my $obj = $self->{headbuff};
        delete $self->{headbuff};
        return $obj if $obj->{component} == "consensus";
    }
    my $fd = $self->{fd};
    while(<$fd>) {
        next unless /^\{/;
        my $obj = decode_json($_);
        return $obj if $obj->{component} == "consensus";
    }
}
