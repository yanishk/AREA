import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "exploreAll.dart";
import "exploreWorkflow.dart";
import "exploreServices.dart";

class TopBarExplore extends StatelessWidget {
  TopBarExplore({
    super.key,
    required this.currentIndex,
  });

  final int currentIndex;
  final List<Widget> _pages = [
    ExplorePageAll(),
    ExplorePageWorkflow(),
    ExplorePageServices(),
  ];

  Future<void>  changePage(int index, BuildContext context) async {
    if (index != currentIndex) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => _pages[index]),
        );
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Row (
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        TextButton(
          onPressed: () {changePage(0, context);},
          child: AutoSizeText("All", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {changePage(1, context);},
          child: AutoSizeText("Workflows", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {changePage(2, context);},
          child: AutoSizeText("Services", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),
      ],
    );
  }
}
