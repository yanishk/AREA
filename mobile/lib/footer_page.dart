import 'package:flutter/material.dart';

import "explore/explorePage.dart";
import "editAccount/editAccount.dart";
import "createWorkflows/create.dart";
import "myWorkflows/displayAllWorkflows.dart";
import 'globalData.dart';

class FooterNavigationBar extends StatelessWidget {
  FooterNavigationBar({
    super.key,
    required this.currentIndex,
  });

  final int currentIndex;
  final List<Widget> _pages = [
    DisplayWorkflows(),
    ExplorePage(),
    CreateWorkflowsPage(),
    EditAccountPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar(
      type: BottomNavigationBarType.fixed,
      showUnselectedLabels: true,
      selectedItemColor: Colors.black,
      unselectedItemColor: Colors.black45,
      backgroundColor: Colors.white,
      selectedLabelStyle:
          const TextStyle(fontSize: 15, fontWeight: FontWeight.bold),
      unselectedLabelStyle:
          const TextStyle(fontSize: 13, fontWeight: FontWeight.bold),
      currentIndex: currentIndex,
      iconSize: 28,
      onTap: (int index) {
        if (index != currentIndex) {
          Navigator.push(
            Globaldata.myContext,
            MaterialPageRoute(builder: (context) => _pages[index]),
          );
        }
      },
      items: const <BottomNavigationBarItem>[
        BottomNavigationBarItem(
          icon: Icon(Icons.call_to_action_outlined),
          label: 'My workflows',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.search),
          label: 'Explore',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.add_circle_outline_outlined),
          label: 'Create',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.person_2_outlined),
          label: 'Profile',
        ),
      ],
    );
  }
}
